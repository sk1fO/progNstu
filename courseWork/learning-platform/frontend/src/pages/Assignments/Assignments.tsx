import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import CodeEditor from '../../Components/editor/CodeEditor';
import { assignmentsAPI, coursesAPI } from '../../services/api';
import { Assignment, Lesson } from '../../types';

const Assignments: React.FC = () => {
  const { lessonId } = useParams<{ lessonId: string }>();
  const [lessons, setLessons] = useState<Lesson[]>([]);
  const [assignments, setAssignments] = useState<Assignment[]>([]);
  const [selectedAssignment, setSelectedAssignment] = useState<Assignment | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      if (!lessonId) return;
      
      try {
        const [lessonsResponse, assignmentsResponse] = await Promise.all([
          coursesAPI.getLessons(Number(lessonId)),
          assignmentsAPI.getAssignments(Number(lessonId))
        ]);
        setLessons(lessonsResponse.data);
        setAssignments(assignmentsResponse.data);
      } catch (error) {
        console.error('Ошибка загрузки данных:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [lessonId]);

  if (loading) return <div>Загрузка...</div>;

  return (
    <div style={{ padding: '20px', display: 'flex', gap: '20px', height: 'calc(100vh - 80px)' }}>
      {/* Список заданий */}
      <div style={{ width: '300px', borderRight: '1px solid #ddd', paddingRight: '20px' }}>
        <h3>Уроки и задания</h3>
        {lessons.map(lesson => (
          <div key={lesson.id} style={{ marginBottom: '20px' }}>
            <h4>{lesson.title}</h4>
            {assignments
              .filter(a => a.lesson_id === lesson.id)
              .map(assignment => (
                <div
                  key={assignment.id}
                  onClick={() => setSelectedAssignment(assignment)}
                  style={{
                    padding: '10px',
                    margin: '5px 0',
                    cursor: 'pointer',
                    backgroundColor: selectedAssignment?.id === assignment.id ? '#f0f0f0' : 'transparent',
                    borderRadius: '4px'
                  }}
                >
                  {assignment.title}
                </div>
              ))}
          </div>
        ))}
      </div>

      {/* Редактор кода */}
      <div style={{ flex: 1 }}>
        {selectedAssignment ? (
          <CodeEditor assignment={selectedAssignment} />
        ) : (
          <div style={{ 
            display: 'flex', 
            alignItems: 'center', 
            justifyContent: 'center', 
            height: '100%',
            color: '#666'
          }}>
            Выберите задание для начала работы
          </div>
        )}
      </div>
    </div>
  );
};

export default Assignments;